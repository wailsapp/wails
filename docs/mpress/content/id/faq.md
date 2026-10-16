---
title: "Pertanyaan Umum"
description: "Jawaban atas pertanyaan paling umum tentang pembuatan aplikasi dengan Wails v3"
slug: "faq"
sourcePath: "faq.md"
---

## Umum

### Apa itu Wails?

Wails adalah framework untuk membuat aplikasi desktop dengan Go dan teknologi web. Logika aplikasi ditulis dalam Go, antarmuka dibuat dengan HTML, CSS, dan JavaScript (atau framework frontend apa pun), lalu Wails merendernya dalam webview native sistem operasi. Hasilnya adalah aplikasi kecil dan cepat yang terasa native: tanpa browser yang dibundel, penggunaan memori rendah, dan satu file biner yang biasanya berukuran sekitar 10 MB.

### Platform apa saja yang didukung Wails?

| Platform | Persyaratan |
| --- | --- |
| Windows | AMD64 dan ARM64. Menggunakan [runtime WebView2](https://developer.microsoft.com/microsoft-edge/webview2/). |
| macOS | 10.15+ pada Intel (aplikasi dapat menargetkan 10.13+), 11.0+ pada Apple Silicon. Biner universal didukung. |
| Linux | AMD64 dan ARM64. Stack default-nya adalah GTK4 dengan WebKitGTK 6.0 (Ubuntu 24.04+, Debian 13+, Fedora 40+, dan yang serupa). Distribusi yang hanya menyediakan WebKit2GTK 4.1, seperti Ubuntu 22.04, Debian 12, dan RHEL 9, didukung melalui build `-tags gtk3` lama (tersedia hingga v3.1). Distribusi yang hanya memiliki WebKit2GTK 4.0 tidak didukung. Lihat [panduan build Linux](/guides/build/linux/). |
| iOS dan Android | Eksperimental. Lihat [panduan seluler](/guides/mobile/). |

Anda juga dapat menyajikan aplikasi sebagai aplikasi web biasa menggunakan [build server](/guides/server-build/).

Jalankan `wails3 doctor` kapan saja untuk memeriksa sistem Anda dan mendapatkan petunjuk instalasi khusus platform.

### Apa yang saya perlukan untuk memulai?

- Go 1.25 atau yang lebih baru
- Node.js dan npm (untuk build frontend)
- Toolchain platform: WebView2 pada Windows (sudah terinstal pada 10/11), Xcode Command Line Tools pada macOS, serta `gcc` dan paket pengembangan GTK/WebKit pada Linux

`wails3 doctor` memeriksa semua ini untuk Anda dan memberi tahu dengan tepat apa yang belum tersedia. Lihat [Instalasi](/quick-start/installation/) untuk panduan lengkapnya.

### Apakah Wails v3 siap digunakan dalam produksi?

Wails v3 adalah perangkat lunak beta dengan API desktop yang stabil. Aplikasi yang menggunakannya sudah berjalan dalam produksi, tetapi Anda sebaiknya mengujinya secara menyeluruh sebelum melakukan deployment selagi kami menyelesaikan penyempurnaan akhir untuk 3.0. Lihat [halaman status proyek](/status/) untuk mengetahui kondisi terkini. Wails v2 adalah rilis stabil saat ini dan terus menerima perbaikan.

## Pengembangan

### Apakah saya perlu menguasai Go?

Pengetahuan dasar tentang Go akan membantu, tetapi Anda tidak perlu menjadi ahli. Logika aplikasi berada dalam metode Go biasa, dan [tutorial](/tutorials/overview/) akan memandu Anda memahami hal-hal lainnya. Banyak pengembang mempelajari Go sambil membuat aplikasi Wails pertama mereka.

### Bisakah saya menggunakan framework frontend favorit saya?

Ya. Jika dapat di-build menjadi HTML, CSS, dan JavaScript, framework tersebut dapat digunakan dengan Wails. Tersedia templat untuk React, Vue, Svelte, dan JavaScript standar (masing-masing dengan varian TypeScript), sedangkan pilihan lainnya dapat diintegrasikan dalam hitungan menit. Lihat [Framework Frontend](/guides/dev/frontend-frameworks/).

### Bagaimana cara memanggil fungsi Go dari JavaScript?

Daftarkan layanan, lalu Wails akan membuat binding bertipe untuk layanan tersebut:

```go
// Go
type GreetService struct{}

func (g *GreetService) Greet(name string) string {
    return "Hello " + name
}
```

```javascript
// JavaScript
import { GreetService } from "./bindings/changeme";

const message = await GreetService.Greet("World");
```

Binding dibuat ulang secara otomatis selama `wails3 dev`, atau sesuai kebutuhan dengan `wails3 generate bindings`. Lihat [Layanan](/features/bindings/services/).

### Bisakah saya menggunakan TypeScript?

Ya. Generator binding menghasilkan definisi TypeScript untuk layanan Anda beserta tipenya, sehingga pemanggilan ke Go sepenuhnya bertipe.

### Bagaimana cara mengirim peristiwa antara Go dan JavaScript?

```go
// Go
app.Event.Emit("time", time.Now().Format(time.RFC1123))
```

```javascript
// JavaScript
import { Events } from "@wailsio/runtime";

Events.On("time", (event) => {
    console.log(event.data);
});
```

Nama peristiwa harus sama persis. Lihat [Referensi Peristiwa](/guides/events-reference/).

### Bagaimana cara men-debug aplikasi saya?

Jalankan `wails3 dev`, lalu klik kanan di dalam jendela untuk membuka alat pengembang browser, sama seperti saat Anda mengembangkan aplikasi web. Server pengembangan juga mendukung hot reload untuk frontend Anda. Lihat [Debugging](/guides/dev/debugging/).

## Build & Distribusi

### Bagaimana cara membuat build untuk produksi?

```bash
wails3 build
```

File biner Anda akan ditempatkan di `bin/`. Build produksi sudah menerapkan default yang sesuai (tag build, `-trimpath`, simbol yang dihapus), sehingga tidak diperlukan flag tambahan untuk menghasilkan file biner yang ringkas.

### Bisakah saya melakukan kompilasi silang?

Bisa, dengan batasan tertentu. Kompilasi silang Go murni tidak berlaku karena setiap platform menggunakan pustaka webview native, tetapi tersedia dukungan yang baik untuk kasus umum:

```bash
# Different architecture, same OS
wails3 build GOOS=windows GOARCH=arm64

# macOS universal binary
wails3 task darwin:build:universal
```

Pembuatan build untuk Linux dari sistem operasi lain menggunakan toolchain berbasis Docker. Lihat [Build Lintas Platform](/guides/build/cross-platform/) untuk matriks lengkapnya.

### Bagaimana cara membuat installer atau paket?

```bash
wails3 package
```

Proses ini menghasilkan format native platform, dan [panduan Installer](/guides/installers/) membahas NSIS pada Windows, bundel `.app` dan DMG pada macOS, serta paket Linux.

### Bagaimana cara menandatangani kode aplikasi saya?

Penandatanganan pada Windows dan macOS (termasuk notarisasi) dibahas langkah demi langkah dalam [panduan Penandatanganan](/guides/build/signing/).

## Fitur

### Bisakah saya membuat beberapa jendela?

Ya, dukungan multi-jendela tersedia secara native di v3:

```go
window1 := app.Window.New()
window2 := app.Window.New()
```

Lihat [Beberapa Jendela](/features/windows/multiple/).

### Apakah Wails mendukung baki sistem?

Ya, termasuk menu dan handler klik:

```go
systemTray := app.SystemTray.New()
systemTray.SetIcon(iconBytes)
systemTray.SetMenu(myMenu)
```

Lihat [Baki Sistem](/features/menus/systray/).

### Bisakah saya menggunakan dialog native?

Ya. Dialog berkas, dialog pesan, dan dialog pertanyaan semuanya menggunakan implementasi native:

```go
path, err := app.Dialog.OpenFile().
    SetTitle("Select File").
    PromptForSingleSelection()
```

Lihat [Dialog](/features/dialogs/overview/).

### Apakah Wails mendukung pembaruan otomatis?

Ya. Wails v3 menyertakan pembaru mandiri bawaan (`app.Updater`) dengan penyedia yang dapat dipasang untuk GitHub Releases, keygen.sh, dan Sparkle AppCast, verifikasi tanda tangan kriptografis, serta UI bawaan yang dapat Anda sesuaikan temanya atau ganti. Lihat panduan [Pembaruan Dalam Aplikasi](/guides/updater/) dan tutorial [Aplikasi Wails yang Memperbarui Diri](/tutorials/04-self-update-a-wails-app/).

## Pemecahan Masalah

### Ada sesuatu yang tidak berfungsi. Dari mana saya harus mulai?

```bash
wails3 doctor
```

Fitur ini memverifikasi toolchain Anda, mencantumkan dependensi yang belum tersedia beserta perintah instalasinya, dan mencetak informasi versi yang sebaiknya Anda sertakan dalam setiap laporan bug.

### Build saya gagal

Solusi yang biasanya berhasil, secara berurutan:

1. `go mod tidy`
2. `cd frontend && npm install` (`node_modules` yang tidak tersedia merupakan penyebab paling umum)
3. Perbarui CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
4. Di Linux, periksa `wails3 doctor` untuk mengetahui paket GTK/WebKit yang belum tersedia

### Binding saya tidak tersedia atau sudah usang

```bash
wails3 generate bindings
```

Binding dibuat ulang secara otomatis dalam mode pengembangan; jika Anda menambahkan layanan baru atau mengubah signature metode di luar `wails3 dev`, buat ulang binding secara manual.

### Event tidak dipicu

Nama event harus sama persis antara `app.Event.Emit("name", ...)` di Go dan `Events.On("name", ...)` di JavaScript. Periksa dahulu kesalahan ketik dan perbedaan huruf besar-kecil.

### Saya menemukan bug

Silakan [buka issue](https://github.com/wailsapp/wails/issues) dan sertakan output `wails3 doctor` Anda. [Panduan umpan balik](/feedback/) menjelaskan cara membuat laporan yang mudah ditindaklanjuti.

## Bermigrasi dari v2

### Haruskah saya bermigrasi dari v2 ke v3?

v3 menghadirkan dukungan multi-jendela, API berbasis layanan yang lebih rapi, pembaru bawaan, sistem build yang jauh lebih fleksibel, dan performa yang lebih baik. Proyek baru sebaiknya dimulai dengan v3. Untuk proyek yang sudah ada, [Panduan Migrasi](/migration/v2-to-v3/) menjelaskan perbedaannya langkah demi langkah.

### Apakah v2 akan tetap dipelihara?

Ya. v2 akan terus menerima perbaikan sementara v3 bergerak menuju rilis stabilnya.

### Bisakah saya menjalankan v2 dan v3 secara berdampingan?

Ya. CLI keduanya merupakan biner terpisah (`wails` dan `wails3`), dan modulnya memiliki jalur impor yang berbeda, sehingga proyek dengan versi mayor yang berbeda dapat berjalan berdampingan tanpa masalah pada satu mesin.

## Komunitas

### Bagaimana cara mendapatkan bantuan?

- [Discord](https://discord.gg/JDdSxwjhGf) untuk pertanyaan singkat dan diskusi
- [GitHub Discussions](https://github.com/wailsapp/wails/discussions) untuk pertanyaan yang memerlukan pembahasan lebih panjang
- [GitHub Issues](https://github.com/wailsapp/wails/issues) untuk bug

### Bagaimana cara berkontribusi?

Lihat [Panduan Kontribusi](/contributing/). Perbaikan bug selalu diterima; fungsionalitas baru dan perubahan perilaku publik menggunakan PR draf [WEP (Wails Enhancement Proposal)](https://github.com/wailsapp/wails/blob/master/v3/wep/README.md). Diskusi informal di Discord atau GitHub Discussions bersifat opsional.

### Di mana saya dapat menemukan contoh?

Repositori ini menyertakan lebih dari 60 contoh yang dapat dijalankan, mencakup jendela, dialog, event, baki sistem, layanan, dan lainnya: [v3/examples](https://github.com/wailsapp/wails/tree/master/v3/examples).

## Masih Ada Pertanyaan?

Tanyakan di [Discord](https://discord.gg/JDdSxwjhGf) atau [buka diskusi](https://github.com/wailsapp/wails/discussions).
