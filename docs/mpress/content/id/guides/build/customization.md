---
title: "Penyesuaian Build"
description: "Sesuaikan proses build Anda menggunakan Task dan Taskfile.yml"
slug: "guides/build/customization"
sourcePath: "guides/build/customization.md"
---

## Ikhtisar

Sistem build Wails adalah alat fleksibel dan andal yang dirancang untuk menyederhanakan proses build aplikasi Wails Anda. Sistem ini memanfaatkan [Task](https://taskfile.dev), sebuah task runner yang memungkinkan Anda mendefinisikan dan menjalankan tugas dengan mudah. Meskipun sistem build v3 digunakan secara default, Wails mendorong pendekatan "gunakan alat Anda sendiri", sehingga pengembang dapat menyesuaikan proses build sesuai kebutuhan.

Pelajari lebih lanjut cara menggunakan Task dalam [dokumentasi resmi](https://taskfile.dev/usage/).

## Task: Inti Sistem Build

[Task](https://taskfile.dev) adalah alternatif modern untuk Make yang ditulis dalam Go. Task menggunakan berkas YAML untuk mendefinisikan tugas dan dependensinya. Dalam sistem build Wails, [Task](https://taskfile.dev) berperan penting dalam mengatur proses build.

`Taskfile.yml` utama terletak di root proyek, sedangkan tugas khusus platform didefinisikan dalam berkas `build/<platform>/Taskfile.yml`. Berkas `Taskfile.yml` umum dalam direktori `build` berisi tugas bersama yang digunakan di berbagai platform.

@filetree

- Project Root
  - Taskfile.yml
  - build
    - windows/Taskfile.yml
    - darwin/Taskfile.yml
    - linux/Taskfile.yml
    - Taskfile.yml
@end

## Taskfile.yml

Berkas `Taskfile.yml` di root proyek merupakan titik masuk utama sistem build. Berkas ini mendefinisikan tugas beserta dependensinya. Berikut adalah berkas `Taskfile.yml` default:

```yaml
version: '3'

includes:
  common: ./build/Taskfile.yml
  windows: ./build/windows/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
  linux: ./build/linux/Taskfile.yml

vars:
  APP_NAME: "myproject"
  BIN_DIR: "bin"
  VITE_PORT: '{{.WAILS_VITE_PORT | default 9245}}'

tasks:
  build:
    summary: Builds the application
    cmds:
      - task: "{{OS}}:build"

  package:
    summary: Packages a production build of the application
    cmds:
      - task: "{{OS}}:package"

  run:
    summary: Runs the application
    cmds:
      - task: "{{OS}}:run"

  dev:
    summary: Runs the application in development mode
    cmds:
      - wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}


```

## Taskfile Khusus Platform

Setiap platform memiliki Taskfile sendiri yang terletak di direktori platform di bawah direktori `build`. Berkas-berkas ini mendefinisikan tugas inti untuk platform tersebut. Setiap Taskfile menyertakan tugas bersama dari berkas `build/Taskfile.yml`.

### Windows

Lokasi: `build/windows/Taskfile.yml`

Taskfile khusus Windows mencakup tugas untuk melakukan build, mengemas, dan menjalankan aplikasi di Windows. Fitur utamanya meliputi:

- Melakukan build dengan flag produksi opsional
- Menghasilkan berkas ikon `.ico`
- Menghasilkan berkas `.syso` Windows
- Membuat penginstal NSIS untuk pengemasan

### Linux

Lokasi: `build/linux/Taskfile.yml`

Taskfile khusus Linux mencakup tugas untuk melakukan build, mengemas, dan menjalankan aplikasi di Linux. Fitur utamanya meliputi:

- Melakukan build dengan flag produksi opsional
- Membuat paket AppImage, deb, rpm, dan Arch Linux
- Menghasilkan berkas `.desktop` untuk aplikasi Linux

### macOS

Lokasi: `build/darwin/Taskfile.yml`

Taskfile khusus macOS mencakup tugas untuk melakukan build, mengemas, dan menjalankan aplikasi di macOS. Fitur utamanya meliputi:

- Melakukan build biner untuk arsitektur amd64, arm64, dan universal (keduanya)
- Menghasilkan berkas ikon `.icns`
- Membuat bundel `.app` untuk distribusi
- Menandatangani bundel `.app` secara ad hoc
- Mengatur flag build dan variabel lingkungan khusus macOS

## Eksekusi Tugas dan Alias Perintah

Perintah `wails3 task` merupakan versi tertanam dari [Taskfile](https://taskfile.dev) yang menjalankan tugas-tugas yang didefinisikan dalam `Taskfile.yml` Anda.

Perintah `wails3 build` dan `wails3 package` masing-masing merupakan alias untuk `wails3 task build` dan `wails3 task package`. Saat Anda menjalankan perintah ini, Wails secara internal menerjemahkannya menjadi eksekusi tugas yang sesuai:

- `wails3 build` → `wails3 task build`
- `wails3 package` → `wails3 task package`

### Meneruskan Parameter ke Tugas

Anda dapat meneruskan variabel CLI ke tugas menggunakan format `KEY=VALUE`. Variabel ini diteruskan melalui perintah alias:

```bash
# These are equivalent:
wails3 build PLATFORM=linux CONFIG=production
wails3 task build PLATFORM=linux CONFIG=production

# Package with custom version:
wails3 package VERSION=2.0.0 OUTPUT=myapp.pkg
```

Dalam `Taskfile.yml`, Anda dapat mengakses variabel tersebut menggunakan sintaks templat Go:

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - go build -tags {{.CONFIG | default "debug"}} -o myapp
```

### Variabel penyesuaian build Go

Proyek yang dihasilkan menyediakan variabel aditif untuk tag build Go, flag linker, dan CGO. Nilai tag harus berupa nama yang dipisahkan koma, tanpa opsi `-tags`; Wails menggabungkannya dengan tag yang diperlukan oleh mode build dan platform yang dipilih.

| Variabel | Berlaku untuk |
| --- | --- |
| `APP_TAGS` | Semua build yang dijalankan melalui Taskfile |
| `APP_TAGS_LINUX` | Build Linux |
| `APP_TAGS_DARWIN` | Build macOS |
| `APP_TAGS_WINDOWS` | Build Windows |
| `APP_TAGS_ANDROID` | Build Android |
| `APP_TAGS_IOS` | Build iOS |
| `APP_TAGS_SERVER` | Build mode server |
| `APP_LDFLAGS` | Flag linker tambahan untuk semua build melalui Taskfile |
| `APP_CGO_ENABLED` | Mengganti pengaturan CGO dengan `0` atau `1` untuk build desktop dan server |
| `EXTRA_TAGS` | Tag tambahan untuk satu pemanggilan |

Contohnya:

```bash
wails3 build APP_TAGS=sqlite_fts5,netgo APP_TAGS_LINUX=myapp_linux
wails3 build EXTRA_TAGS=diagnostics
APP_LDFLAGS='-X example.com/myapp/internal/version.Value=1.2.3' wails3 build
```

Tetapkan nilai `APP_*` di `Taskfile.yml` utama jika ingin menjadikannya nilai default proyek yang persisten. `KEY=value` yang diberikan bersama perintah memiliki prioritas tertinggi dan mengganti variabel tersebut untuk pemanggilan itu. Nilai literal di Taskfile utama memiliki prioritas lebih tinggi daripada variabel lingkungan proses dengan nama yang sama. Jika ekspresi yang merujuk dirinya sendiri dengan nilai default seperti di bawah tetap dipertahankan, lingkungan proses digunakan saat tidak ada nilai dari baris perintah. Gunakan `EXTRA_TAGS` untuk menambahkan tag bagi satu build tanpa mengganti nilai `APP_TAGS` yang persisten.

Saat `APP_CGO_ENABLED` kosong, Linux dan macOS menggunakan default `1`, Windows menggunakan `0`, build server natif mempertahankan default Go pada host, dan build server Docker menggunakan `0`. Build Android dan iOS selalu menggunakan CGO dan mempertahankan `CGO_ENABLED=1`; `APP_CGO_ENABLED` tidak mengganti pengaturan toolchain seluler tersebut. Kompilasi silang desktop dan build server berbasis Docker menerima nilai `APP_*` yang sesuai, sama seperti build natif melalui Taskfile.

@note{type="info" title="Proyek yang sudah ada"}
`Taskfile.yml` utama dimiliki oleh proyek, sehingga `wails3 update build-assets` tidak menimpanya. Proyek yang dibuat sebelum variabel ini diperkenalkan harus menambahkan entri berikut secara manual ke blok `vars` utama:

```yaml
vars:
  APP_TAGS: '{{.APP_TAGS | default ""}}'
  APP_TAGS_LINUX: '{{.APP_TAGS_LINUX | default ""}}'
  APP_TAGS_DARWIN: '{{.APP_TAGS_DARWIN | default ""}}'
  APP_TAGS_WINDOWS: '{{.APP_TAGS_WINDOWS | default ""}}'
  APP_TAGS_ANDROID: '{{.APP_TAGS_ANDROID | default ""}}'
  APP_TAGS_IOS: '{{.APP_TAGS_IOS | default ""}}'
  APP_TAGS_SERVER: '{{.APP_TAGS_SERVER | default ""}}'
  APP_LDFLAGS: '{{.APP_LDFLAGS | default ""}}'
  APP_CGO_ENABLED: '{{.APP_CGO_ENABLED | default ""}}'
```

Ganti ekspresi yang merujuk dirinya sendiri dengan nilai default menjadi nilai literal untuk menjadikannya default proyek yang persisten.

@end

Variabel ini digunakan oleh build yang dijalankan melalui Taskfile. Tahap build dalam proyek Xcode yang dihasilkan memanggil Go secara langsung dan berada di luar alur penyesuaian Taskfile ini, sehingga build yang dijalankan dari Xcode saat ini tidak menggunakan variabel tersebut.

## Proses Build Umum

Di semua platform, proses build biasanya mencakup langkah-langkah berikut:

1. Merapikan modul Go
2. Melakukan build frontend
3. Menghasilkan ikon
4. Mengompilasi kode Go dengan flag khusus platform
5. Mengemas aplikasi (khusus platform)

## Menyesuaikan Proses Build

Meskipun sistem build v3 menyediakan konfigurasi default yang andal, Anda dapat dengan mudah menyesuaikannya dengan kebutuhan proyek. Dengan mengubah `Taskfile.yml` dan Taskfile khusus platform, Anda dapat:

- Menambahkan tugas baru
- Mengubah tugas yang sudah ada
- Mengubah urutan eksekusi tugas
- Mengintegrasikan alat dan skrip lain

Fleksibilitas ini memungkinkan Anda menyesuaikan proses build dengan kebutuhan spesifik sekaligus tetap memperoleh manfaat dari struktur yang disediakan oleh sistem build Wails.

@note{type="tip" title="Mempelajari Taskfile"}
Kami sangat menyarankan Anda membaca dokumentasi [Taskfile](https://taskfile.dev) untuk memahami cara menggunakan Taskfile secara efektif. Anda dapat mengetahui versi Taskfile yang tertanam dalam Wails CLI dengan menjalankan `wails3 task --version`.

@end

## Mode Pengembangan

Sistem build Wails menyertakan mode pengembangan andal yang meningkatkan pengalaman pengembang dengan menyediakan pemuatan ulang langsung dan penggantian modul secara hot. Mode ini diaktifkan menggunakan perintah `wails3 dev`.

### Cara Kerjanya

Saat Anda menjalankan `wails3 dev`, proses berikut akan berlangsung:

1. Perintah tersebut memeriksa port yang tersedia dan menggunakan 9245 secara default jika port tidak ditentukan.
2. Perintah tersebut menyiapkan variabel lingkungan untuk server pengembangan frontend (Vite).
3. Perintah tersebut memulai pemantau berkas menggunakan pustaka [refresh](https://github.com/atterpac/refresh).

Pustaka [refresh](https://github.com/atterpac/refresh) bertugas memantau perubahan berkas dan memicu build ulang. Pustaka ini menggunakan konfigurasi yang ditetapkan pada kunci `dev_mode` dalam berkas `./build/config.yml`. Pustaka ini dapat dikonfigurasi untuk mengabaikan direktori dan berkas tertentu, menentukan berkas yang akan dipantau, serta menentukan tindakan yang harus dilakukan ketika perubahan terdeteksi. Konfigurasi default sudah berfungsi cukup baik, tetapi Anda dapat menyesuaikannya sesuai kebutuhan.

### Konfigurasi

Berikut adalah contoh strukturnya:

```yaml
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir:
      - .git
      - node_modules
      - frontend
      - bin
    file:
      - .DS_Store
      - .gitignore
      - .gitkeep
    watched_extension:
      - "*.go"
    git_ignore: true
  executes:
    - cmd: wails3 task common:install:frontend:deps
      type: once
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: go mod tidy
      type: blocking
    - cmd: wails3 task build
      type: blocking
    - cmd: wails3 task run
      type: primary
```

Berkas konfigurasi ini memungkinkan Anda untuk:

- Menetapkan jalur akar untuk pemantauan berkas
- Mengonfigurasi tingkat pencatatan log
- Menetapkan waktu debounce untuk peristiwa perubahan berkas
- Mengabaikan direktori, berkas, atau ekstensi berkas tertentu
- Menentukan perintah yang akan dijalankan saat berkas berubah

### Menyesuaikan Mode Pengembangan

Anda dapat menyesuaikan pengalaman mode pengembangan dengan mengubah nilai-nilai ini dalam berkas `config.yml`.

Beberapa penyesuaian yang dapat dilakukan meliputi:

1. Mengubah direktori atau berkas yang dipantau
2. Menyesuaikan waktu debounce untuk mengatur seberapa cepat sistem merespons perubahan
3. Menambahkan atau mengubah perintah eksekusi agar sesuai dengan kebutuhan proyek Anda

### Menggunakan browser untuk pengembangan

Meskipun Wails v2 sepenuhnya mendukung penggunaan browser untuk pengembangan, hal tersebut menimbulkan banyak kebingungan. Aplikasi yang berfungsi di browser belum tentu berfungsi sebagai aplikasi desktop karena tidak semua API browser tersedia di webview.

Untuk pekerjaan pengembangan yang berfokus pada UI, di v3 Anda tetap dapat menggunakan browser dengan mengakses URL Vite di `http://localhost:9245` dalam mode pengembangan. Dengan demikian, Anda dapat menggunakan alat pengembangan browser yang canggih saat mengerjakan gaya dan tata letak. Perlu diketahui bahwa binding Go *tidak akan berfungsi* dalam mode ini. Saat siap menguji fungsionalitas seperti binding dan peristiwa, cukup beralih ke tampilan desktop untuk memastikan semuanya berfungsi sempurna di lingkungan produksi.
