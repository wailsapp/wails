---
title: "Aplikasi Pertama Anda"
description: "Buat aplikasi desktop Wails pertama Anda langkah demi langkah"
slug: "getting-started/your-first-app"
sourcePath: "getting-started/your-first-app.md"
---

Panduan ini menunjukkan cara membuat aplikasi Wails v3 pertama Anda, mulai dari penyiapan proyek dan build hingga alur kerja pengembangan.

<br/>

<br/>

@steps
### Membuat Proyek Baru
Buka terminal Anda, lalu jalankan perintah berikut untuk membuat proyek Wails baru:

```bash
wails3 init -n myfirstapp
```

Perintah ini membuat direktori baru bernama `myfirstapp` yang berisi semua file yang diperlukan.

   <video src="/assets/wails_init.mp4" controls></video>

### Menjelajahi Struktur Proyek
Buka direktori `myfirstapp`. Anda akan menemukan beberapa file dan folder:

@filetree
- build/           Berisi file yang digunakan oleh proses build
  - appicon.png  Ikon aplikasi
  - config.yml   Konfigurasi build
  - Taskfile.yml Build tasks
  - darwin/      File build khusus macOS
    - Info.dev.plist Development configuration
    - Info.plist    Konfigurasi produksi
    - Taskfile.yml  Tugas build macOS
    - icons.icns    Ikon aplikasi macOS
  - linux/       File build khusus Linux
    - Taskfile.yml  Tugas build Linux
    - appimage/     Pemaketan AppImage
      - build.sh  Skrip build AppImage
    - nfpm/        Pemaketan NFPM
      - nfpm.yaml Package configuration
      - scripts/  Skrip build
  - windows/     File build khusus Windows
    - Taskfile.yml        Tugas build Windows
    - icon.ico           Ikon aplikasi Windows
    - info.json          Metadata aplikasi
    - wails.exe.manifest Windows manifest file
    - nsis/              File penginstal NSIS
      - project.nsi                    File proyek NSIS
      - wails_tools.nsh               Skrip pembantu NSIS
- frontend/        File aplikasi frontend
  - index.html   File HTML utama
  - main.js      File JavaScript utama
  - package.json NPM package configuration
  - public/      Aset statis
  - Inter Font License.txt Font license
- .gitignore      File pengabaian Git
- README.md       Dokumentasi proyek
- Taskfile.yml    Tugas proyek
- go.mod          File modul Go
- go.sum          Checksum modul Go
- greetservice.go Greeting service
- main.go         Kode aplikasi utama
@end

Luangkan waktu sejenak untuk menjelajahi file-file ini dan mengenali strukturnya.

@note{type="info"}
Meskipun Wails v3 menggunakan [Task](https://taskfile.dev/) sebagai sistem build default, Anda tetap dapat menggunakan `make` atau sistem build alternatif lainnya.

@end

### Membuat Build Aplikasi Anda
Untuk membuat build aplikasi Anda, jalankan:

```bash
wails3 build
```

Perintah ini mengompilasi versi debug aplikasi Anda dan menyimpannya dalam direktori `bin` baru.

@note{type="info"}
`wails3 build` adalah bentuk singkat dari `wails3 task build` dan akan menjalankan tugas `build` dalam `Taskfile.yml`.

@end

     <video src="/assets/wails_build.mp4" controls></video>

Setelah build selesai, Anda dapat menjalankannya seperti aplikasi biasa lainnya:

@tabs{sync-key="platform"}
[Mac]
```sh
./bin/myfirstapp
```

[Windows]
```sh
bin\myfirstapp.exe
```

[Linux]
```sh
./bin/myfirstapp
```

@end

Anda akan melihat antarmuka pengguna sederhana yang menjadi titik awal aplikasi Anda. Karena ini adalah versi debug, Anda juga akan melihat log di jendela konsol. Log tersebut berguna untuk proses debug.

### Mode Pengembangan
Anda juga dapat menjalankan aplikasi dalam mode pengembangan. Mode ini memungkinkan Anda mengubah kode frontend dan langsung melihat perubahan tersebut pada aplikasi yang sedang berjalan tanpa harus membuat ulang seluruh build aplikasi.

1. Buka jendela terminal baru.
2. Jalankan `wails3 dev`. Aplikasi akan dikompilasi dan dijalankan dalam mode debug.
3. Buka `frontend/index.html` di editor pilihan Anda.
4. Edit kode tersebut dan ubah `Please enter your name below` menjadi `Please enter your name below!!!`.
5. Simpan file tersebut.

Perubahan ini akan segera terlihat pada aplikasi Anda.

Setiap perubahan pada kode backend akan memicu build ulang:

1. Buka `greetservice.go`.
2. Ubah baris yang berisi `return "Hello " + name + "!"` menjadi `return "Hello there " + name + "!"`.
3. Simpan file tersebut.

Aplikasi akan diperbarui dalam hitungan detik.

     <video src="/assets/wails_dev.mp4" controls></video>

### Memaketkan Aplikasi Anda
Setelah aplikasi siap didistribusikan, Anda dapat membuat paket khusus platform:

@tabs{sync-key="platform"}
[Mac]
Untuk membuat bundel `.app`:

```bash
wails3 package
```

Tindakan ini akan membuat build produksi dan memaketkannya menjadi bundel `.app` dalam direktori `bin`.

[Windows]
Untuk membuat penginstal NSIS:

```bash
wails3 package
```

Tindakan ini akan membuat build produksi dan memaketkannya menjadi penginstal NSIS dalam direktori `bin`.

[Linux]
Wails mendukung beberapa format paket untuk distribusi Linux:

```bash
# Create all package types (AppImage, deb, rpm, and Arch Linux)
wails3 package

# Or create specific package types
wails3 task linux:create:appimage  # AppImage format
wails3 task linux:create:deb       # Debian package
wails3 task linux:create:rpm       # Red Hat package
wails3 task linux:create:aur       # Arch Linux package
```

@end

Untuk informasi lebih terperinci tentang opsi dan konfigurasi pemaketan, lihat [Panduan Build dan Pemaketan](/guides/build/building/) kami.

### Menyiapkan Kontrol Versi dan Nama Modul
Proyek Anda dibuat dengan nama modul sementara `changeme`. Sebaiknya perbarui nama ini agar sesuai dengan URL repositori Anda:

1. Buat repositori baru di GitHub (atau layanan hosting Git pilihan Anda)
2. Inisialisasi Git dalam direktori proyek Anda:
  ```bash
  git init
  git add .
  git commit -m "Initial commit"
  ```

3. Tetapkan repositori remote Anda (ganti dengan URL repositori Anda):
  ```bash
  git remote add origin https://github.com/username/myfirstapp.git
  ```

4. Perbarui nama modul dalam `go.mod` agar sesuai dengan URL repositori Anda:
  ```bash
  go mod edit -module github.com/username/myfirstapp
  ```

5. Push kode Anda:
  ```bash
  git push -u origin main
  ```


Langkah ini memastikan nama modul Go Anda selaras dengan konvensi penamaan modul Go dan memudahkan Anda membagikan kode.

@note{type="tip" title="Kiat Profesional"}
Anda dapat mengotomatiskan semua langkah inisialisasi dengan menggunakan flag `-git` saat membuat proyek:

```bash
wails3 init -n myfirstapp -git github.com/username/myfirstapp
```

Fitur ini mendukung berbagai format URL Git:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` atau `ssh://git@github.com/username/project`
- Protokol Git: `git://github.com/username/project`
- Sistem file: `file:///path/to/project.git`

@end

@end

## Selamat!

Anda baru saja membuat, mengembangkan, dan memaketkan aplikasi Wails pertama Anda. Ini hanyalah awal dari berbagai hal yang dapat Anda capai dengan Wails v3.

## Langkah Berikutnya

Jika Anda baru mengenal Wails, kami menyarankan Anda membaca Tutorial kami selanjutnya, yang akan menjadi panduan praktis untuk mempelajari berbagai fitur Wails. Tutorial pertama adalah [Membuat Layanan](/tutorials/01-creating-a-service/).

Jika Anda adalah pengguna yang lebih mahir, lihat [Panduan Build dan Pemaketan](/guides/build/building/) untuk mendapatkan informasi yang lebih mendetail tentang cara menggunakan Wails.
