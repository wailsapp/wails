---
title: "Flag dan Contoh CLI"
description: "Cari flag dan contoh penggunaan untuk setiap perintah CLI Wails"
slug: "guides/cli"
sourcePath: "guides/cli.md"
---

Gunakan referensi ini jika Anda sudah mengetahui perintah yang diperlukan. Untuk menelusuri perintah berdasarkan tujuannya, mulailah dari [Ikhtisar Perintah CLI](/reference/cli/).

CLI Wails menyediakan serangkaian perintah lengkap untuk membantu Anda mengembangkan, membangun, dan memelihara aplikasi Wails.

## Perintah Inti

Perintah inti adalah perintah utama yang digunakan untuk membuat, mengembangkan, dan membangun proyek.

Semua perintah CLI menggunakan format berikut: `wails3 <command>`.

### `init`

Menginisialisasi proyek Wails baru. Selama inisialisasi ini, perintah `go mod tidy` dijalankan untuk memperbarui paket proyek. Langkah ini dapat dilewati dengan menggunakan flag `-skipgomodtidy` bersama perintah `init`.

```bash
wails3 init [flags]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-p` | Nama paket Go | `main` |
| `-t` | Nama atau URL templat | `vanilla` |
| `-n` | Nama proyek |  |
| `-d` | Direktori proyek | `.` |
| `-q` | Sembunyikan output | `false` |
| `-l` | Cantumkan templat | `false` |
| `-mod` | Path modul Go (dihitung dari `-git` jika tidak ditentukan) |  |
| `-git` | URL repositori Git |  |
| `-s` | Lewati peringatan saat menggunakan templat jarak jauh | `false` |
| `-productname` | Nama produk | `My Product` |
| `-productdescription` | Deskripsi produk | `My Product Description` |
| `-productversion` | Versi produk | `0.1.0` |
| `-productcompany` | Nama perusahaan | `My Company` |
| `-productcopyright` | Pemberitahuan hak cipta | `© now, My Company` |
| `-productcomments` | Komentar file | `This is a comment` |
| `-productidentifier` | Pengidentifikasi produk |  |
| `-skipgomodtidy` | Lewati go mod tidy | `false` |

Flag `-git` menerima berbagai format URL Git:

- HTTPS: `https://github.com/username/project`
- SSH: `git@github.com:username/project` atau `ssh://git@github.com/username/project`
- Protokol Git: `git://github.com/username/project`
- Sistem file: `file:///path/to/project.git`

Jika diberikan, flag ini akan:

1. Menginisialisasi repositori Git di direktori proyek
2. Menetapkan URL yang ditentukan sebagai remote origin
3. Memperbarui nama modul di `go.mod` agar sesuai dengan URL repositori
4. Menambahkan semua file

### `dev`

Menjalankan aplikasi dalam mode pengembangan. Anda akan mendapatkan tampilan langsung kode frontend sehingga dapat membuat perubahan dan langsung melihat hasilnya pada aplikasi yang sedang berjalan tanpa perlu membangun ulang seluruh aplikasi. Perubahan pada kode Go Anda juga akan terdeteksi, lalu aplikasi akan dibangun ulang dan dijalankan kembali secara otomatis.

```bash
wails3 dev [flags]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-config` | Path file konfigurasi | `./build/config.yml` |
| `-port` | Port server pengembangan Vite | `9245` |
| `-s` | Aktifkan HTTPS | `false` |

@note{type="info"}
Ini setara dengan menjalankan `wails3 task dev` dan akan menjalankan tugas `dev` dalam Taskfile utama proyek. Anda dapat menyesuaikannya dengan mengedit file `Taskfile.yml`.

@end

### `build`

Membangun versi produksi menggunakan Taskfile yang dihasilkan. Gunakan `wails3 dev` untuk hot reload, atau berikan `DEV=true` untuk build pengembangan. Secara default, build dibuat untuk platform dan arsitektur saat ini.

```bash
wails3 build [flags] [CLI variables...]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-tags` | Tag build Go tambahan (dipisahkan dengan koma) |  |

Anda dapat memberikan variabel CLI untuk menyesuaikan build:

```bash
wails3 build GOOS=linux
```

Gunakan flag `-tags` untuk memberikan tag build Go khusus:

```bash
# Build with legacy GTK3 + WebKit2GTK 4.1 on Linux (default is GTK4 + WebKitGTK 6.0)
wails3 build -tags gtk3

# Build in server mode (no GUI)
wails3 build -tags server

# Multiple tags
wails3 build -tags gtk3,customtag
```

Tag diteruskan sebagai `EXTRA_TAGS` ke Taskfile yang mendasarinya.

@note{type="info"}
Ini setara dengan menjalankan `wails3 task build`, yang akan menjalankan tugas `build` dalam Taskfile utama proyek. Semua variabel CLI yang diberikan kepada `build` diteruskan ke tugas yang mendasarinya. Anda dapat menyesuaikan proses build dengan mengedit file `Taskfile.yml`.

@end

### `package`

Membuat paket khusus platform untuk distribusi.

```bash
wails3 package [CLI variables...]
```

Tugas pemaketan bawaan menerima `ARCH` untuk memilih arsitektur target:

```bash
wails3 package ARCH=amd64
```

`VERSION` dan `OUTPUT` hanya memengaruhi Taskfile khusus yang secara eksplisit membaca variabel tersebut; keduanya tidak mengganti versi paket atau nama file bawaan.

#### Jenis Paket

Jenis paket berikut tersedia untuk setiap platform:

| Platform | Jenis Paket |
| --- | --- |
| Windows | `.exe` |
| macOS | `.app`, |
| Linux | `.AppImage`, `.deb`, `.rpm`, `.archlinux` |

@note{type="info"}
Ini setara dengan `wails3 task package`, yang menjalankan tugas `package` dalam Taskfile utama proyek. Semua variabel CLI yang diteruskan ke `package` akan diteruskan ke tugas yang mendasarinya. Anda dapat menyesuaikan proses pemaketan dengan mengedit berkas `Taskfile.yml`.

@end

Tugas build bawaan membaca `GOOS`, `ARCH`, dan `DEV`; `PLATFORM` dan `CONFIG` dalam tugas khusus di bawah ini merupakan contoh yang ditentukan pengguna.

### `task`

Menjalankan tugas yang ditentukan dalam Taskfile.yml proyek Anda. Ini adalah versi tersemat dari [Taskfile](https://taskfile.dev) yang memungkinkan Anda menentukan dan menjalankan tugas build, pengujian, dan deployment khusus.

```bash
wails3 task [taskname] [CLI variables...] [flags]
```

#### Variabel CLI

Anda dapat meneruskan variabel ke tugas dengan format `KEY=VALUE`:

```bash
wails3 task build GOOS=linux
wails3 task deploy ENV=staging VERSION=1.2.3
```

Variabel ini dapat diakses dalam Taskfile.yml menggunakan sintaks templat Go:

```yaml
tasks:
  build:
    cmds:
      - echo "Building for {{.PLATFORM | default "darwin"}}"
      - echo "Config: {{.CONFIG | default "debug"}}"
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-h` | Menampilkan cara penggunaan Task | `false` |
| `-i` | Membuat Taskfile.yml baru | `false` |
| `-list` | Mencantumkan tugas beserta deskripsinya | `false` |
| `-list-all` | Mencantumkan semua tugas (dengan atau tanpa deskripsi) | `false` |
| `-json` | Memformat daftar tugas sebagai JSON | `false` |
| `-status` | Keluar dengan kode bukan nol jika tugas belum mutakhir | `false` |
| `-f` | Memaksa eksekusi meskipun tugas sudah mutakhir | `false` |
| `-w` | Mengaktifkan mode pemantauan untuk tugas yang ditentukan | `false` |
| `-v` | Mengaktifkan mode verbose | `false` |
| `-version` | Menampilkan versi Task | `false` |
| `-s` | Menonaktifkan echo | `false` |
| `-p` | Menjalankan tugas secara paralel | `false` |
| `-dry` | Mengompilasi dan menampilkan tugas tanpa menjalankannya | `false` |
| `-summary` | Menampilkan ringkasan tentang suatu tugas | `false` |
| `-x` | Meneruskan kode keluar tugas | `false` |
| `-dir` | Menetapkan direktori eksekusi |  |
| `-taskfile` | Memilih Taskfile yang akan dijalankan |  |
| `-output` | Menetapkan gaya output: [interleaved|group|prefixed] |  |
| `-c` | Output berwarna (diaktifkan secara bawaan) | `true` |
| `-C` | Membatasi jumlah tugas yang dijalankan secara bersamaan |  |
| `-interval` | Interval pemantauan perubahan (dalam detik) |  |

#### Contoh

```bash
# Run the default task
wails3 task

# Run a specific task
wails3 task test

# Run a task with variables
wails3 task build GOOS=windows ARCH=amd64

# List all available tasks
wails3 task --list

# Run multiple tasks in parallel
wails3 task -p task1 task2 task3

# Watch for changes and re-run task
wails3 task -w dev
```

### `mcp`

Memulai server MCP proyek Wails untuk pengelolaan proyek dengan bantuan agen. Server ini terpisah dari server MCP yang dikompilasi ke dalam aplikasi yang sedang berjalan: `wails3 mcp` mengelola berkas proyek dan perintah siklus hidup, sedangkan server MCP aplikasi mengendalikan WebView yang sedang berjalan.

```bash
wails3 mcp [flags]
```

Transport dipilih secara otomatis:

- Saat host MCP menjalankan Wails dengan stdin/stdout yang disalurkan melalui pipe, server menggunakan **stdio**.
- Saat dijalankan secara interaktif di terminal, server menggunakan **Streamable HTTP** pada `127.0.0.1` dan meminta port yang tersedia dari sistem operasi.

Gunakan `--stdio` atau `--http` untuk memilih transport secara eksplisit. Gunakan `--port 0` untuk memilih port loopback yang tersedia dalam mode HTTP.

#### Flag MCP

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `--root` | Root proyek yang diizinkan. Path dan symlink di luar root tersebut ditolak. | Direktori saat ini |
| `--token` | Token sesi/bearer untuk alat yang melakukan perubahan dan mengendalikan proses. Jika tidak tersedia, menggunakan `WAILS_MCP_TOKEN`. | Dihasilkan secara aman |
| `--stdio` | Memaksa penggunaan transport stdio. | Otomatis |
| `--http` | Memaksa penggunaan transport Streamable HTTP. | Otomatis |
| `--port` | Port HTTP; `0` memilih port loopback yang tersedia. | `0` |

Dalam mode HTTP, Wails mencetak endpoint dan token bearer ke stderr. Dalam mode stdio, token disertakan dalam petunjuk inisialisasi MCP. Server tidak menyediakan eksekusi perintah shell secara arbitrer. Templat jarak jauh dan remote Git memerlukan persetujuan eksplisit melalui input `allowExternal` milik alat tersebut.

### `doctor`

Melakukan pemeriksaan sistem dan menampilkan laporan status.

```bash
wails3 doctor
```

## Perintah Pembuatan

Perintah pembuatan membantu membuat berbagai aset proyek, seperti binding, ikon, dan berkas build. Semua perintah pembuatan menggunakan perintah dasar: `wails3 generate <command>`.

### `generate bindings`

Menghasilkan binding dan model untuk kode Go Anda.

```bash
wails3 generate bindings [flags] [patterns...]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-f` | Flag build Go tambahan |  |
| `-d` | Direktori keluaran | `frontend/bindings` |
| `-models` | Nama berkas model | `models` |
| `-index` | Nama berkas indeks | `index` |
| `-ts` | Hasilkan TypeScript | `false` |
| `-i` | Gunakan antarmuka TS | `false` |
| `-b` | Gunakan runtime yang dibundel | `false` |
| `-names` | Gunakan nama sebagai pengganti ID | `false` |
| `-noindex` | Lewati berkas indeks | `false` |
| `-noevents` | Lewati pembuatan binding terkait peristiwa | `false` |
| `-dry` | Simulasi tanpa perubahan | `false` |
| `-silent` | Mode senyap | `false` |
| `-v` | Keluaran debug | `false` |
| `-clean` | Bersihkan direktori keluaran sebelum pembuatan | `true` |

### `generate build-assets`

Menghasilkan aset build untuk aplikasi Anda.

```bash
wails3 generate build-assets [flags]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-name` | Nama proyek |  |
| `-dir` | Direktori keluaran | `build` |
| `-silent` | Sembunyikan keluaran | `false` |
| `-company` | Nama perusahaan |  |
| `-productname` | Nama produk |  |
| `-description` | Deskripsi produk |  |
| `-version` | Versi produk |  |
| `-identifier` | Pengidentifikasi produk | `com.wails.[name]` |
| `-copyright` | Pemberitahuan hak cipta |  |
| `-comments` | Komentar berkas |  |

### `generate icons`

Menghasilkan ikon aplikasi.

```bash
wails3 generate icons [flags]
```

#### Flag

| Flag | Deskripsi | Bawaan |
| --- | --- | --- |
| `-input` | Berkas PNG masukan | Wajib |
| `-windowsfilename` | Nama berkas keluaran Windows |  |
| `-macfilename` | Nama berkas keluaran macOS |  |
| `-sizes` | Ukuran ikon (dipisahkan dengan koma) | `256,128,64,48,32,16` |
| `-example` | Hasilkan ikon contoh | `false` |
| `-iconcomposerinput` | Berkas Icon Composer masukan (`.icon`) |  |
| `-macassetdir` | Direktori keluaran untuk aset Mac (Assets.car + icns) |  |

#### Icon Composer (macOS)

Di macOS 26+, Anda dapat menggunakan berkas `.icon` Icon Composer untuk menghasilkan `Assets.car` dan `icons.icns`:

```bash
wails3 generate icons -iconcomposerinput build/appicon.icon -macassetdir build
```

Proses ini mengompilasi berkas `.icon` menggunakan perintah `actool` dari Apple. Memerlukan Xcode dengan `actool` versi 26 atau yang lebih baru.

Saat menggunakan Icon Composer, atur `cfBundleIconName` dalam `build/config.yml` agar sesuai dengan nama berkas `.icon` (tanpa ekstensi):

```yaml
info:
  cfBundleIconName: "appicon"
```

Jika tidak diatur dan `Assets.car` tersedia, nilai bawaannya adalah `"appicon"`.

### `generate syso`

Menghasilkan file .syso Windows.

```bash
wails3 generate syso [flags]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-manifest` | Path ke file manifes | Wajib |
| `-icon` | Path ke file ikon | Wajib |
| `-info` | Path ke file informasi versi |  |
| `-arch` | Arsitektur target | GOARCH saat ini |
| `-out` | Nama file keluaran | `rsrc_windows_[arch].syso` |

### `generate .desktop`

Menghasilkan file .desktop Linux.

```bash
wails3 generate .desktop [flags]
```

#### Flag

| Flag | Deskripsi | Bawaan |
| --- | --- | --- |
| `-name` | Nama aplikasi | Wajib |
| `-exec` | Path ke berkas yang dapat dieksekusi | Wajib |
| `-icon` | Path ikon |  |
| `-categories` | Kategori aplikasi | `Utility` |
| `-comment` | Komentar aplikasi |  |
| `-terminal` | Jalankan di terminal | `false` |
| `-keywords` | Kata kunci pencarian |  |
| `-version` | Versi aplikasi |  |
| `-genericname` | Nama generik |  |
| `-startupnotify` | Tampilkan notifikasi saat dimulai | `false` |
| `-mimetype` | Tipe MIME yang didukung |  |
| `-output` | Nama file keluaran | `[name].desktop` |

### `generate runtime`

Menghasilkan versi runtime yang telah dikompilasi sebelumnya.

```bash
wails3 generate runtime
```

### `generate constants`

Menghasilkan konstanta JavaScript dari kode Go.

```bash
wails3 generate constants
```

### `generate webview2bootstrapper`

Menghasilkan penginstal bootstrap WebView2 untuk Windows yang siap didistribusikan.

```bash
wails3 generate webview2bootstrapper [flags]
```

### `generate template`

Membuat kerangka direktori templat proyek baru.

```bash
wails3 generate template [flags]
```

### `generate appimage`

Menghasilkan AppImage Linux.

```bash
wails3 generate appimage [flags]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-binary` | Path ke file biner | Wajib |
| `-icon` | Path ke file ikon | Wajib |
| `-desktop` | Path ke file .desktop | Wajib |
| `-builddir` | Direktori build | Direktori sementara |
| `-output` | Direktori keluaran | `.` |

## Perintah Layanan

Perintah layanan membantu mengelola layanan Wails. Semua perintah layanan menggunakan perintah dasar: `wails3 service <command>`.

### `service init`

Menginisialisasi layanan baru.

```bash
wails3 service init [flags]
```

#### Flag

| Flag | Deskripsi | Bawaan |
| --- | --- | --- |
| `-n` | Nama layanan | `example_service` |
| `-d` | Deskripsi layanan | `Example service` |
| `-p` | Nama paket |  |
| `-o` | Direktori keluaran | `.` |
| `-q` | Sembunyikan keluaran | `false` |
| `-a` | Nama penulis |  |
| `-v` | Versi |  |
| `-w` | URL situs web |  |
| `-r` | URL repositori |  |
| `-l` | Lisensi |  |

## Perintah Alat

Perintah alat menyediakan utilitas untuk pengembangan dan penelusuran kesalahan. Semua perintah alat menggunakan perintah dasar: `wails3 tool <command>`.

### `tool checkport`

Memeriksa apakah suatu port terbuka. Berguna untuk menguji apakah vite sedang berjalan.

```bash
wails3 tool checkport [flags]
```

#### Flag

| Flag | Deskripsi | Bawaan |
| --- | --- | --- |
| `-port` | Port yang akan diperiksa | `9245` |
| `-host` | Host yang akan diperiksa | `localhost` |

### `tool watcher`

Memantau file dan menjalankan perintah ketika file berubah.

```bash
wails3 tool watcher [flags]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-config` | Jalur file konfigurasi | `./build/config.yml` |
| `-ignore` | Pola yang akan diabaikan |  |
| `-include` | Pola yang akan disertakan |  |

### `tool cp`

Menyalin file.

```bash
wails3 tool cp
```

### `tool buildinfo`

Menampilkan informasi build aplikasi.

```bash
wails3 tool buildinfo
```

### `tool version`

Menaikkan versi semantik berdasarkan flag yang diberikan.

```bash
wails3 tool version [flags]
```

#### Flag

| Opsi | Deskripsi | Default |
| --- | --- | --- |
| `-v` | Versi saat ini yang akan dinaikkan |  |
| `-major` | Naikkan versi mayor | `false` |
| `-minor` | Naikkan versi minor | `false` |
| `-patch` | Naikkan versi patch | `false` |
| `-prerelease` | Naikkan versi prarilis (misalnya, dari alpha.5 menjadi alpha.6) | `false` |

Perintah mengikuti urutan prioritas: mayor > minor > patch > prarilis. Perintah ini mempertahankan prefiks "v" jika ada dalam versi masukan, serta semua komponen prarilis dan metadata.

Contoh penggunaan:

```bash
wails3 tool version -v 1.2.3 -major      # Output: 2.0.0
wails3 tool version -v v1.2.3 -minor     # Output: v1.3.0
wails3 tool version -v 1.2.3-alpha -patch # Output: 1.2.4-alpha
wails3 tool version -v v3.1.0-alpha.5 -prerelease # Output: v3.1.0-alpha.6
```

### `tool package`

Menghasilkan paket Linux (deb, rpm, archlinux).

```bash
wails3 tool package [flags]
```

#### Opsi

| Flag | Deskripsi | Bawaan |
| --- | --- | --- |
| `-format` | Format paket (deb, rpm, archlinux) | `deb` |
| `-name` | Nama file yang dapat dieksekusi | `myapp` |
| `-config` | Jalur file konfigurasi |  |
| `-out` | Direktori keluaran | `.` |

### `tool lipo`

Membuat biner universal macOS dengan menggabungkan biner khusus arsitektur.

```bash
wails3 tool lipo [flags]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-output` | Jalur biner keluaran |  |

### `tool capabilities`

Memeriksa kemampuan build sistem (ketersediaan GTK4/GTK3 di Linux).

```bash
wails3 tool capabilities
```

### `tool docker-mounts`

Menghasilkan flag pemasangan volume Docker untuk kompilasi silang. Menghasilkan flag `-v` untuk cache modul Go dan semua direktif `replace` lokal dalam `go.mod`, untuk digunakan dalam perintah `docker run` Taskfile.

```bash
wails3 tool docker-mounts
```

### `tool has`

Memeriksa apakah alat atau kemampuan tersedia, lalu mencetak `true` atau `false` ke stdout. Dirancang untuk digunakan dalam variabel `sh:` Taskfile sebagai alternatif lintas platform untuk `command -v`.

Gunakan `|` untuk memeriksa apakah salah satu dari beberapa alternatif tersedia.

```bash
wails3 tool has <tool>
```

#### Contoh

```bash
# Check for a C compiler (gcc or clang)
wails3 tool has gcc|clang

# Check for a specific tool
wails3 tool has git
wails3 tool has node
```

#### Penggunaan dalam Taskfile

```yaml
vars:
  HAS_CC:
    sh: 'wails3 tool has gcc|clang'
```

### `tool has-cc`

@note{type="caution" title="Tidak digunakan lagi"}
`wails3 tool has-cc` tidak digunakan lagi. Perbarui Taskfile Anda agar menggunakan `wails3 tool has gcc|clang` sebagai gantinya.

@end

Alias yang kompatibel dengan versi sebelumnya untuk `wails3 tool has gcc|clang`. Memeriksa apakah `gcc` atau `clang` tersedia dalam PATH, lalu mencetak `true` atau `false`.

```bash
wails3 tool has-cc
```

## Perintah Pembaruan

Perintah pembaruan membantu mengelola dan memperbarui aset proyek. Semua perintah pembaruan menggunakan perintah dasar: `wails3 update <command>`.

### `update cli`

Memperbarui Wails CLI ke versi baru.

```bash
wails3 update cli [flags]
```

#### Flag

| Flag | Deskripsi | Bawaan |
| --- | --- | --- |
| `-pre` | Perbarui ke prarilis terbaru | `false` |
| `-version` | Perbarui ke versi tertentu |  |
| `-nocolour` | Nonaktifkan keluaran berwarna | `false` |

Perintah update cli memungkinkan Anda memperbarui instalasi Wails CLI. Secara default, perintah ini memperbarui ke rilis stabil terbaru. Anda dapat menggunakan flag `-pre` untuk memperbarui ke versi prarilis terbaru, atau menentukan versi tertentu dengan flag `-version`.

Setelah memperbarui, jangan lupa memperbarui file go.mod proyek Anda agar menggunakan versi yang sama:

```bash
require github.com/wailsapp/wails/v3 v3.x.x
```

### `update build-assets`

Memperbarui aset build menggunakan file konfigurasi yang diberikan.

```bash
wails3 update build-assets [flags]
```

#### Flag

| Flag | Deskripsi | Bawaan |
| --- | --- | --- |
| `-config` | Jalur file konfigurasi |  |
| `-dir` | Direktori keluaran | `build` |
| `-silent` | Sembunyikan keluaran | `false` |
| `-company` | Nama perusahaan |  |
| `-productname` | Nama produk |  |
| `-description` | Deskripsi produk |  |
| `-version` | Versi produk |  |
| `-identifier` | Pengidentifikasi produk |  |
| `-copyright` | Pemberitahuan hak cipta |  |
| `-comments` | Komentar file |  |

## Perintah Utilitas

Perintah utilitas menyediakan pintasan praktis untuk tugas umum. Gunakan perintah ini secara langsung dengan perintah dasar: `wails3 <command>`.

### `docs`

Membuka dokumentasi Wails di peramban bawaan Anda.

```bash
wails3 docs
```

### `releasenotes`

Menampilkan catatan rilis untuk versi saat ini atau versi yang ditentukan.

```bash
wails3 releasenotes [flags]
```

#### Flag

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-v` | Versi yang catatan rilisnya akan ditampilkan |  |
| `-n` | Nonaktifkan keluaran berwarna | `false` |

### `version`

Mencetak versi Wails saat ini.

```bash
wails3 version
```

### `sponsor`

Membuka halaman sponsor Wails di peramban bawaan Anda.

```bash
wails3 sponsor

```
