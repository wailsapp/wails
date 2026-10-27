---
title: "Menulis dokumentasi dengan M-Press"
description: "Unduh dan instal M-Press, pratinjau dokumentasi Wails, lalu kirim kontribusi dokumentasi."
sourcePath: "contributing/documentation.md"
---

Wails v3 menggunakan **M-Press** untuk mengubah Markdown menjadi situs dokumentasi ini. Anda dapat memperbaiki kesalahan ketik, menyempurnakan penjelasan, menambahkan contoh atau panduan, maupun memperbaiki terjemahan. Perbaikan yang hanya menyangkut dokumentasi tidak memerlukan issue sebelumnya atau pengujian kode yang gagal.

Anda memerlukan Git untuk menyiapkan kontribusi dan akun GitHub untuk mengirim pull request. Berkas M-Press siap eksekusi yang dapat diunduh mencakup generator dokumentasi dan perender diagram; Anda tidak memerlukan Go, Node.js, Wails CLI, atau dependensi build desktop untuk melihat pratinjau dokumentasi. Python 3 diperlukan untuk skrip validasi tambahan milik repositori.

## Mulai dari halaman dokumentasi

Buka **Kontribusi → Dokumentasi** di bilah sisi untuk membuka panduan ini. Untuk mengerjakan halaman tertentu, pilih **ikon edit** di header halaman tersebut, lalu pilih **Edit halaman ini secara lokal**. Salin perintah untuk sistem operasi Anda dan jalankan di terminal.

Perintah tersebut menggunakan M-Press jika sudah terinstal; jika belum, penginstalnya akan mengunduh rilis resmi dan memverifikasi checksum. Perintah itu menyiapkan working copy kontribusi lokal dan membuka halaman yang dipilih di M-Press. Sebelum menjalankan pemeriksaan repositori, cocokkan `mpress version` dengan versi yang dipatok di bawah karena penginstal otomatis dapat mengunduh rilis yang lebih baru.

Anda dapat mengedit sumber di editor sendiri sementara pratinjau lokal dibuat ulang. Wizard juga menyediakan tindakan terpisah untuk **Tinjau perubahan**, **Commit perubahan**, **Push cabang kontribusi**, dan **Buka draf pull request**. Ikuti panduan pemeriksaan dan PR di bawah sebelum mengirimkannya. Pembuatan terjemahan bersifat opsional dan mungkin memerlukan penyedia yang telah dikonfigurasi; pengeditan dan pratinjau Markdown tidak memerlukannya.

Untuk penyiapan manual yang dapat direproduksi, ikuti langkah-langkah berikut.

## Unduh dan instal M-Press

Skrip build repositori mengunduh **rilis stabil terbaru M-Press**. Gunakan rilis yang sama secara lokal: pilih metode instalasi di bawah ini, jalankan `mpress version`, lalu bandingkan hasilnya dengan versi yang ditampilkan di halaman rilis.

### Unduh berkas siap eksekusi

Buka [rilis stabil terbaru M-Press](https://github.com/leaanthony/mpress/releases/latest), atau unduh arsip yang sesuai di bawah ini. ARM64 mencakup Apple Silicon; AMD64 adalah arsitektur yang umum digunakan pada PC Intel/AMD.

| Komputer | AMD64 / Intel | ARM64 / Apple Silicon |
| --- | --- | --- |
| macOS | [mpress-darwin-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-amd64.tar.gz) | [mpress-darwin-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-darwin-arm64.tar.gz) |
| Linux | [mpress-linux-amd64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-amd64.tar.gz) | [mpress-linux-arm64.tar.gz](https://github.com/leaanthony/mpress/releases/latest/download/mpress-linux-arm64.tar.gz) |
| Windows | [mpress-windows-amd64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-amd64.zip) | [mpress-windows-arm64.zip](https://github.com/leaanthony/mpress/releases/latest/download/mpress-windows-arm64.zip) |

Unduh [checksums.txt](https://github.com/leaanthony/mpress/releases/latest/download/checksums.txt) dari rilis yang sama. Di direktori unduhan Anda, hitung hash SHA-256 arsip tersebut, lalu bandingkan dengan baris untuk nama file itu di `checksums.txt`:

@tabs
[macOS]
```sh
shasum -a 256 mpress-darwin-arm64.tar.gz
```

[Linux]
```sh
sha256sum mpress-linux-amd64.tar.gz
```

[Windows PowerShell]
```powershell
Get-FileHash .\mpress-windows-amd64.zip -Algorithm SHA256
```
@end

Gunakan nama berkas unduhan Anda yang sebenarnya jika memilih arsitektur lain. Setelah hash cocok, ekstrak dan instal berkas siap eksekusi:

@tabs
[macOS]
Dari direktori unduhan, menggunakan arsip Apple Silicon:

```sh
mkdir -p mpress-download
tar -xzf mpress-darwin-arm64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Untuk Mac berbasis Intel, gunakan `mpress-darwin-amd64.tar.gz` sebagai gantinya. Tambahkan baris `export PATH` ke `~/.zshrc` agar perintah tetap tersedia di terminal baru.

[Linux]
Dari direktori unduhan:

```sh
mkdir -p mpress-download
tar -xzf mpress-linux-amd64.tar.gz -C mpress-download
mkdir -p "$HOME/.local/bin"
install -m 755 mpress-download/mpress "$HOME/.local/bin/mpress"
export PATH="$HOME/.local/bin:$PATH"
mpress version
```

Untuk ARM64, gunakan `mpress-linux-arm64.tar.gz` sebagai gantinya. Tambahkan baris `export PATH` ke berkas startup shell Anda, misalnya `~/.bashrc`, agar perintah tetap tersedia di terminal baru.

[Windows PowerShell]
Dari direktori unduhan:

```powershell
Expand-Archive .\mpress-windows-amd64.zip -DestinationPath .\mpress-download
New-Item -ItemType Directory -Force "$HOME\bin"
Copy-Item .\mpress-download\mpress.exe "$HOME\bin\mpress.exe"
$env:Path = "$HOME\bin;$env:Path"
mpress version
```

Untuk ARM64, gunakan `mpress-windows-arm64.zip` sebagai gantinya. Tambahkan direktori `bin` pengguna Anda ke **Path** pengguna melalui **Variabel Lingkungan** Windows, lalu buka kembali terminal agar perintah tetap tersedia pada sesi berikutnya.
@end

### Instal dengan Go

Jika Anda sudah menggunakan Go, Anda dapat membangun rilis stabil terbaru sebagai alternatif. Periksa `go.mod` rilis tersebut untuk mengetahui toolchain yang diperlukan; pemilihan toolchain otomatis Go mungkin akan mengunduhnya. Rilis yang digunakan untuk proses dokumentasi ini memerlukan Go 1.26.10 atau yang lebih baru.

```sh
go install github.com/leaanthony/mpress/cmd/mpress@latest
mpress version
```

Go menginstal perintah di `GOBIN`, atau di direktori `bin` dalam `go env GOPATH` jika `GOBIN` tidak ditetapkan. Tambahkan direktori tersebut ke `PATH` Anda dan buka kembali terminal jika `mpress` tidak ditemukan.

## Dapatkan sumber dokumentasi

[Fork wailsapp/wails](https://github.com/wailsapp/wails/fork), lalu jalankan perintah berikut dengan mengganti `YOUR_USERNAME` dengan nama pengguna GitHub Anda:

```sh
git clone https://github.com/YOUR_USERNAME/wails.git
cd wails
git remote add upstream https://github.com/wailsapp/wails.git
git fetch upstream master
git switch -c docs/my-change upstream/master
```

Langkah ini memulai cabang dokumentasi Anda dari cabang upstream **master**. Gunakan nama cabang deskriptif Anda sendiri sebagai pengganti `docs/my-change`.

## Jalankan pratinjau langsung

Dari root repositori Wails, tempat `mpress.yaml` berada:

```sh
mpress dev
```

Buka [http://127.0.0.1:3000](http://127.0.0.1:3000) di browser Anda. Untuk panduan ini, kunjungi [http://127.0.0.1:3000/contributing/documentation/](http://127.0.0.1:3000/contributing/documentation/). Biarkan perintah tetap berjalan selama Anda mengedit; perubahan sumber yang disimpan akan membuat ulang pratinjau. Tekan **Ctrl+C** untuk menghentikannya sebelum menjalankan perintah validasi di bawah.

Repositori sudah berisi konfigurasi M-Press. Anda tidak perlu menginisialisasi situs baru. Untuk menggunakan port lain, jalankan `mpress dev --port 3001` dan buka alamat localhost yang sesuai.

## Edit atau tambahkan halaman

Halaman sumber berbahasa Inggris berada di `docs/mpress/content/`. Terjemahan berada dalam folder bahasa seperti `fr/` dan `id/`. Petunjuk kontribusi pada ikon edit menunjukkan berkas sumber untuk halaman yang sedang dibuka.

Untuk halaman yang sudah ada, edit Markdown-nya dan pertahankan metadata YAML di bagian atas. Gunakan judul, paragraf, daftar, tautan, dan blok kode berpagar. Pertahankan pasangan komponen M-Press seperti `@note` / `@end` dan `@tabs` / `@end`. [Dokumentasi penulisan M-Press](https://m-press.me/authoring/) menjelaskan sintaks yang didukung.

Untuk panduan baru, buat berkas seperti `docs/mpress/content/guides/my-guide.md`:

```markdown
---
title: "My guide"
description: "The task this guide helps someone complete."
slug: "guides/my-guide"
sourcePath: "guides/my-guide.md"
---

Explain the task, give working steps, and show the expected result.
```

Tambahkan halaman ke bagian yang sesuai dalam `docs/mpress/content/_nav.yaml`, kecuali bagian tersebut sudah menggunakan `autogenerate`. Pastikan pembaca dapat membukanya melalui navigasi dan tautan dari halaman terkait. Halaman baru berbahasa Inggris juga memerlukan salinan lengkap untuk setiap locale agar pemeriksaan cakupan dapat lulus; koordinasikan hal ini dengan pengelola atau kontributor terjemahan.

Output yang dihasilkan berada di `docs/mpress/site/`. Commit halaman sumber serta perubahan navigasi atau aset statis yang dimaksudkan; berkas situs yang dihasilkan dan cache lokal tidak disertakan dalam PR.

## Periksa kontribusi Anda

Dari root repositori:

```sh
python3 docs/mpress/scripts/check_translations.py
mpress build --strict --no-purge-css
mpress check
python3 docs/mpress/scripts/check_site.py docs/mpress/site
```

Gunakan `python` di Windows jika Python 3 diinstal dengan cara tersebut. Jika [Task](https://taskfile.dev/) sudah terinstal, `task docs:build` menjalankan pemeriksaan yang sama menggunakan `mpress` yang tersedia di PATH Anda.

Pemeriksaan terjemahan akan gagal jika ada file terjemahan yang hilang, kosong, atau tidak lagi terkait dengan sumbernya. Bagian terjemahan yang usang atau rusak secara default dilaporkan sebagai peringatan; `--strict` membuat temuan tersebut juga menyebabkan pemeriksaan gagal. Perubahan pada prosa bahasa Inggris yang sudah ada dapat diajukan sementara sinkronisasi terjemahan menyusul. Perbaikan pada satu terjemahan dapat hanya mengubah bahasa tersebut. Pertahankan perintah, nama API, tautan, contoh kode, dan koneksi diagram; jangan gunakan placeholder bahasa Inggris sebagai konten terjemahan.

Baca halaman yang diubah di browser dan jalankan setiap contoh kode yang Anda ubah. Jika pemeriksaan gagal pada halaman yang tidak diubah, sertakan output dan konteks yang relevan dalam PR Anda alih-alih melonggarkan pemeriksaan tersebut.

## Ajukan pull request

Tinjau diff Anda, commit perubahan sumber, lalu push branch kontribusi ke fork Anda. Buka PR terhadap **master** di `wailsapp/wails`. Jelaskan perubahan yang dilakukan, tautkan halaman yang terdampak, dan cantumkan pemeriksaan yang Anda jalankan. Sertakan tangkapan layar sebelum dan sesudah untuk perubahan tata letak serta detail platform/versi untuk contoh kode yang diubah.

Situs produksi bersifat statis; pratinjau dan validasi lokal tidak memerlukan kredensial deployment maupun penyedia terjemahan. Maintainer menangani publikasi setelah peninjauan.

Untuk kontribusi kode dan usulan fitur, lihat [Berkontribusi pada Wails](/contributing/). Untuk internal framework, lihat [Ikhtisar Teknis](/contributing/overview/).
