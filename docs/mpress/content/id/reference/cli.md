---
title: "Referensi CLI"
description: "Referensi lengkap untuk perintah CLI Wails"
slug: "reference/cli"
sourcePath: "reference/cli.md"
---

## Ikhtisar

CLI Wails (`wails3`) adalah titik masuk baris perintah untuk membuat, mengembangkan, membangun, menandatangani, mengemas, dan memeriksa aplikasi Wails 3. Sebagian besar orkestrasi build didelegasikan ke Taskfile masing-masing proyek (di bawah `build/` dalam proyek Anda) — banyak perintah `wails3` hanya berupa pembungkus tipis yang menjalankan tugas tertentu.

Untuk mendapatkan bantuan terbaru tentang perintah apa pun, jalankan:

```bash
wails3 --help
wails3 <command> --help
```

## Siklus hidup proyek

| Perintah | Deskripsi |
| --- | --- |
| `wails3 init` | Buat proyek baru dari templat. Flag: `-n` (nama proyek), `-t` (templat, default `vanilla`), `-p` (nama paket Go, default `main`), `-d` (direktori proyek, default `.`), `-q` (mode senyap), `-l` (tampilkan daftar templat), `-mod` (jalur modul Go), `--git` (URL repositori Git), `--skipgomodtidy`, `-s` (lewati peringatan templat jarak jauh), `--productname`/`--productdescription`/`--productversion`/`--productcompany`/`--productcopyright`/`--productcomments`/`--productidentifier`. |
| `wails3 dev` | Jalankan aplikasi dalam mode pengembangan dengan pemuatan ulang langsung frontend. Flag: `--config` (default `./build/config.yml`), `--port` (port pengembangan Vite), `-s` (aktifkan HTTPS). |
| `wails3 build` | Bangun proyek. Pembungkus tipis untuk tugas Taskfile `build`. Flag: `--tags` (diteruskan sebagai `EXTRA_TAGS=`), `--obfuscated` (bangun dengan Garble; lihat [Build yang Diobfusikasi](/guides/build/obfuscation/)), `--garbleargs` (flag tambahan yang diteruskan ke `garble` sebelum subperintah `build`). |
| `wails3 package` | Jalankan tugas Taskfile `package` khusus platform. |
| `wails3 task [name]` | Jalankan tugas Taskfile apa pun; jika nama tidak diberikan, `--list` menampilkan semua tugas yang terdaftar. |
| `wails3 mcp` | Mulai server MCP proyek. Secara otomatis menggunakan stdio untuk proses yang diluncurkan agen atau Streamable HTTP loopback untuk penggunaan terminal interaktif. |
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
| `generate bindings` | Buat binding dari Go ke frontend. Flag: `-d` (direktori output), `-models`, `-index`, `-ts`, `-i` (antarmuka), `-b` (bundel), `-names` (hasilkan `Call.ByName`), `-noevents`, `-noindex`, `-dry`, `-silent`, `-v`, `-clean` (default `true`), `-f`, `-obfuscated` (buat `wails_obfuscated.gen.go` dengan ID binding yang stabil untuk build Garble; lihat [Build yang Diobfusikasi](/guides/build/obfuscation/)), `-obfuscated-output` (direktori untuk berkas yang dihasilkan; default-nya adalah direktori paket utama). Menerima pola paket (misalnya `./...`); jika tidak ada pola yang diberikan, perintah ini otomatis menggunakan direktori saat ini. |
| `generate icons` | Konversikan PNG sumber ke format ikon platform. Flag: `-input`, `-windowsfilename`, `-macfilename`, `-iconcomposerinput`, `-macassetdir`. |
| `generate build-assets` | Buat isi direktori `build/` (cuplikan Taskfile, berkas NSIS, `Info.plist`, templat `.desktop`, dan sebagainya) dari `build/config.yml`. |
| `generate runtime` | Buat ulang `/wails/runtime.js` bawaan yang dikirimkan ke webview. |
| `generate syso` | Buat berkas sumber daya `.syso` Windows (ikon + manifes + informasi versi). |
| `generate webview2bootstrapper` | Buat penginstal bootstrap WebView2 untuk Windows. |
| `generate constants` | Buat konstanta nama peristiwa JS dari jenis peristiwa Go. |
| `generate template` | Buat kerangka templat proyek baru. |
| `generate .desktop` | Buat berkas `.desktop` Linux (digunakan oleh AppImage/DEB/RPM). |
| `generate appimage` | Buat direktori build AppImage. |

## Pembaruan

`wails3 update <subcommand>`:

| Subperintah | Deskripsi |
| --- | --- |
| `update build-assets` | Perbarui direktori `build/` dari `build/config.yml` (mempertahankan perubahan pengguna jika memungkinkan). |
| `update cli` | Jalankan pembaruan mandiri untuk biner `wails3`. |

## Penandatanganan kode & pengemasan

| Perintah | Deskripsi |
| --- | --- |
| `wails3 setup signing` | Wizard interaktif yang mengonfigurasi penandatanganan untuk platform yang terdeteksi di `build/`. Flag: `--platform` (dapat diulang; secara default dideteksi otomatis dari direktori build). |
| `wails3 setup entitlements` | Wizard interaktif untuk entitlement macOS. Flag: `--output` (path; default `build/darwin/entitlements.plist`). |
| `wails3 sign [GOOS=…]` | Pembungkus yang menjalankan tugas Taskfile `*:sign` khusus platform untuk OS saat ini (atau OS yang ditentukan melalui `GOOS`). |
| `wails3 tool sign` | Titik masuk tingkat rendah untuk penandatanganan langsung. Flag: `--input`, `--output`, `--verbose`, `--certificate`, `--password`, `--thumbprint`, `--timestamp`, `--identity`, `--entitlements`, `--hardened-runtime`, `--notarize`, `--keychain-profile`, `--pgp-key`, `--pgp-password`, `--role`. |

**Tidak ada** subperintah `wails3 signing` — untuk kredensial keychain, gunakan `xcrun notarytool store-credentials`, dan untuk kunci PGP, gunakan `gpg` secara langsung (wizard `wails3 setup signing` mengotomatiskan keduanya).

## Alat

`wails3 tool <subcommand>`:

| Subperintah | Deskripsi |
| --- | --- |
| `tool checkport` | Periksa apakah port TCP terbuka (berguna saat menunggu Vite). |
| `tool watcher` | Jalankan perintah setiap kali file yang dipantau berubah. |
| `tool cp` | Salin file lintas platform. |
| `tool buildinfo` | Tampilkan informasi build Go yang disematkan dari sebuah berkas biner. |
| `tool package` | Buat paket Linux (`deb`, `rpm`, `archlinux`) dari `build/linux/nfpm`. |
| `tool version` | Naikkan versi semantik proyek. |
| `tool lipo` | Gabungkan beberapa berkas biner untuk arsitektur macOS yang berbeda menjadi berkas biner universal. |
| `tool capabilities` | Periksa sistem untuk mengetahui ketersediaan GTK3/GTK4/WebKit. |
| `tool sign` | (Lihat [Penandatanganan kode & pemaketan](#penandatanganan-kode--pengemasan).) |

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
| `ios xcode:gen` | Buat proyek Xcode di direktori output. |

## Path output build

- Berkas biner native ditempatkan di `bin/<APP_NAME>` (atau `bin/<APP_NAME>.exe` di Windows). Tidak ada `build/bin/`.
- Output yang telah dipaketkan (`.app`, `.dmg`, penginstal NSIS, MSIX, DEB/RPM/AppImage) juga ditempatkan di `bin/` (atau di subdirektori khusus platform yang dibuat oleh tugas Taskfile terkait).

## Flag global

| Flag | Berlaku untuk | Deskripsi |
| --- | --- | --- |
| `--no-colour` | Semua perintah | Nonaktifkan warna ANSI pada output CLI. |

---

**Ada pertanyaan?** Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau lihat [contoh](https://github.com/wailsapp/wails/tree/master/v3/examples).
