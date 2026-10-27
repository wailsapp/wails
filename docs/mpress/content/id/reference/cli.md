---
title: "Ikhtisar Perintah CLI"
description: "Temukan perintah CLI Wails berdasarkan tujuan: proyek, build, penandatanganan, alat, dan layanan"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

Mulai dari sini untuk menemukan perintah berdasarkan tujuan. Untuk flag khusus perintah dan contoh penggunaan, lihat [Flag dan Contoh CLI](/guides/cli/).

## Ikhtisar

CLI Wails (`wails3`) adalah titik masuk baris perintah untuk membuat, mengembangkan, membangun, menandatangani, mengemas, dan memeriksa aplikasi Wails 3. Sebagian besar orkestrasi build didelegasikan ke Taskfile per proyek (di bawah `build/` dalam proyek Anda) — banyak perintah `wails3` hanyalah pembungkus tipis yang menjalankan tugas tertentu.

Untuk mendapatkan bantuan terbaru mengenai perintah apa pun, jalankan:

```bash
wails3 --help
wails3 <command> --help
```

## Siklus hidup proyek

| Perintah | Deskripsi |
| --- | --- |
| `wails3 init` | Buat proyek baru dari templat. Flag: `-n` (nama proyek), `-t` (templat, nilai default `vanilla`), `-p` (nama paket Go, nilai default `main`), `-d` (direktori proyek, nilai default `.`), `-q` (mode senyap), `-l` (cantumkan templat), `-mod` (jalur modul Go), `--git` (URL repositori Git), `--skipgomodtidy`, `-s` (lewati peringatan templat jarak jauh), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | Jalankan aplikasi dalam mode pengembangan dengan hot reload frontend. Flag: `--config` (nilai default `./build/config.yml`), `--port` (port pengembangan Vite), `-s` (aktifkan HTTPS). |
| `wails3 build` | Bangun proyek. Pembungkus tipis untuk tugas Taskfile `build`. Flag: `--tags` (diteruskan sebagai `EXTRA_TAGS=`), `--obfuscated` (bangun dengan Garble; lihat [Build yang Diobfusikasi](/guides/build/obfuscation/)), `--garbleargs` (flag tambahan yang diteruskan ke `garble` sebelum subperintah `build`). |
| `wails3 package` | Jalankan tugas Taskfile `package` khusus platform. |
| `wails3 task [name]` | Jalankan tugas Taskfile apa pun; jika nama tidak diberikan, `--list` menampilkan semua tugas yang terdaftar. |
| `wails3 mcp` | Mulai server MCP proyek. Secara otomatis menggunakan stdio untuk proses yang diluncurkan agen atau Streamable HTTP melalui loopback untuk penggunaan terminal interaktif. |
| `wails3 doctor` | Cetak laporan diagnostik lingkungan Anda. |
| `wails3 doctor-ng` | Varian TUI yang lebih baru dari `doctor`. |
| `wails3 version` | Cetak versi CLI. |
| `wails3 releasenotes` | Cetak catatan rilis terbaru. |
| `wails3 docs` | Buka situs dokumentasi di peramban Anda. |
| `wails3 sponsor` | Buka halaman sponsor. |

## Pembuatan

`wails3 generate <subcommand>`:

| Subperintah | Deskripsi |
| --- | --- |
| `generate bindings` | Buat binding dari Go ke frontend. Flag: `-d` (direktori keluaran), `-models`, `-index`, `-ts`, `-i` (antarmuka), `-b` (bundel), `-names` (hasilkan `Call.ByName`), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean` (nilai default `true`), `-f`, `-obfuscated` (buat `wails_obfuscated.gen.go` dengan ID binding stabil untuk build Garble; lihat [Build yang Diobfusikasi](/guides/build/obfuscation/)), `-obfuscated-output` (direktori untuk berkas yang dibuat; nilai default-nya adalah direktori paket utama). Menerima pola paket (misalnya `./...`); jika tidak ada yang diberikan, gunakan direktori saat ini. |
| `generate icons` | Konversikan PNG sumber menjadi format ikon platform. Flag: `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | Buat isi direktori `build/` (cuplikan Taskfile, berkas NSIS, `Info.plist`, templat `.desktop`, dan sebagainya) dari `build/config.yml`. |
| `generate runtime` | Buat ulang `/wails/runtime.js` siap pakai yang disertakan untuk webview. |
| `generate syso` | Buat berkas sumber daya Windows `.syso` (ikon + manifes + informasi versi). |
| `generate webview2bootstrapper` | Buat penginstal bootstrap WebView2 untuk Windows. |
| `generate constants` | Buat konstanta nama peristiwa JS dari tipe peristiwa Go. |
| `generate template` | Buat kerangka templat proyek baru. |
| `generate .desktop` | Buat berkas Linux `.desktop` (digunakan oleh AppImage/DEB/RPM). |
| `generate appimage` | Buat direktori build AppImage. |

## Pembaruan

`wails3 update <subcommand>`:

| Subperintah | Deskripsi |
| --- | --- |
| `update build-assets` | Segarkan direktori `build/` dari `build/config.yml` (mempertahankan perubahan pengguna jika memungkinkan). |
| `update cli` | Perbarui sendiri biner `wails3`. |

## Penandatanganan kode & pengemasan

| Perintah | Deskripsi |
| --- | --- |
| `wails3 signing` | Wizard interaktif yang mengonfigurasi penandatanganan untuk platform yang dideteksinya di `build/`. Flag: `--platform` (dapat diulang; nilai default-nya adalah deteksi otomatis dari direktori build). |
| `wails3 entitlements` | Wizard interaktif untuk entitlement macOS. Flag: `--output` (jalur; nilai default `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | Pembungkus yang menjalankan tugas Taskfile `*:sign` khusus platform untuk OS saat ini (atau OS yang ditentukan melalui `GOOS`). |
| `wails3 tool sign` | Titik masuk penandatanganan langsung tingkat rendah. Flag: `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

Dalam rilis ini, parser CLI menyediakan `signing` dan `entitlements` pada tingkat teratas. Periksa `wails3 signing --help` dan `wails3 entitlements --help`; `wails3 setup signing` tidak memilih wizard penandatanganan. Untuk mengelola kredensial secara langsung, gunakan `xcrun notarytool store-credentials` atau `gpg`.

## Alat

`wails3 tool <subcommand>`:

| Subperintah | Deskripsi |
| --- | --- |
| `tool checkport` | Periksa apakah port TCP terbuka (berguna untuk menunggu Vite). |
| `tool watcher` | Jalankan perintah setiap kali berkas yang dipantau berubah. |
| `tool cp` | Salin berkas lintas platform. |
| `tool buildinfo` | Tampilkan informasi build Go yang disematkan dalam berkas biner. |
| `tool package` | Buat paket Linux (`deb`, `rpm`, `archlinux`) dari `build/linux/nfpm`. |
| `tool version` | Naikkan versi semantik proyek. |
| `tool lipo` | Gabungkan berkas biner dari beberapa arsitektur macOS menjadi satu berkas biner universal. |
| `tool capabilities` | Periksa ketersediaan GTK3/GTK4/WebKit pada sistem. |
| `tool sign` | (Lihat [Penandatanganan kode & pengemasan](#penandatanganan-kode--pengemasan).) |

## Layanan

`wails3 service <subcommand>`:

| Subperintah | Deskripsi |
| --- | --- |
| `service init` | Buat kerangka paket layanan baru. |

## iOS

`wails3 ios <subcommand>`:

| Subperintah | Deskripsi |
| --- | --- |
| `ios overlay:gen` | Buat overlay Go untuk shim jembatan iOS. |
| `ios xcode:gen` | Buat proyek Xcode di direktori keluaran. |

## Lokasi keluaran build

- Berkas biner native ditempatkan di `bin/<APP_NAME>` (atau `bin/<APP_NAME>.exe` di Windows). Tidak ada `build/bin/`.
- Keluaran yang dikemas (`.app`, `.dmg`, penginstal NSIS, MSIX, DEB/RPM/AppImage) juga ditempatkan di `bin/` (atau di subdirektori khusus platform yang dibuat oleh tugas Taskfile terkait).

## Flag bersama

| Flag | Berlaku untuk | Deskripsi |
| --- | --- | --- |
| `-nocolour` | Perintah yang menyediakannya di `--help` (misalnya `init` dan `update cli`) | Nonaktifkan warna ANSI dalam keluaran CLI. Flag ini tidak diterima oleh semua perintah, termasuk `version`. |

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh](https://github.com/wailsapp/wails/tree/releases/v3.0.0/v3/examples).
